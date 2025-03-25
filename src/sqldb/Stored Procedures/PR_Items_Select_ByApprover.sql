CREATE PROCEDURE [dbo].[PR_Items_Select_ByApprover]
    @Approver [VARCHAR](100),
    @RequestType [VARCHAR](100) = NULL,
	@Organization [VARCHAR](100) = NULL,
    @Offset INT = 0,
    @Filter INT = 10
AS
BEGIN
    WITH [ItemResults] AS (
        SELECT
            [dbo].[UidToString]([I].[Id]) AS [Id],
            [I].[Subject] AS [Subject],
            [dbo].[UidToString]([A].[Id]) AS [ApplicationId],
            [A].[Name] AS [ApplicationName],
            [dbo].[UidToString]([AM].[Id]) AS [ApplicationModuleId],
            [AM].[Name] AS [ApplicationModuleName],
            [I].[Created] AS [Created],
            [I].[CreatedBy] AS [RequestedBy],
            [I].[Body] AS [Body],
            (SELECT STRING_AGG([ApproverEmail], ',') FROM [ApprovalRequestApprovers] WHERE [ItemId] = [I].[Id] GROUP BY [ItemId]) AS [Approvers]
        FROM [dbo].[Items] AS [I]
        INNER JOIN [dbo].[ApplicationModules] AS [AM] ON [I].[ApplicationModuleId] = [AM].[Id]
        INNER JOIN [dbo].[Applications] AS [A] ON [AM].[ApplicationId] = [A].[Id]
        INNER JOIN [dbo].[ApprovalRequestApprovers] AS [ARA] ON [I].[Id] = [ARA].[ItemId]
        WHERE
        [I].[IsApproved] IS NULL
        AND [ARA].[ApproverEmail] = @Approver
        AND (
			@RequestType IS NULL OR
			(@RequestType IS NOT NULL AND [I].[ApplicationModuleId] = @RequestType)
		)
        AND (
			@Organization IS NULL OR
			(@Organization IS NOT NULL AND [I].[Body] LIKE '%'+@Organization+'%')
		)
    )

    SELECT *
    INTO [#TemporaryItemResults]
    FROM [ItemResults]

    SELECT 
        * 
    FROM [#TemporaryItemResults]
    ORDER BY [Created] DESC
    OFFSET @Offset ROWS 
    FETCH NEXT @Filter ROWS ONLY

    SELECT COUNT(*) AS Total FROM [#TemporaryItemResults]
END