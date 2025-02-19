CREATE PROCEDURE [dbo].[PR_Items_Select_ByModuleId]
(
	@Search VARCHAR(50) = '',
	@Offset INT = 0,
	@Filter INT = 10,
	@IsApproved int = 4,
	@ModuleId UNIQUEIDENTIFIER
)
AS
BEGIN
	SELECT
		dbo.UidToString([I].[Id]) AS [ItemId],
		dbo.UidToString([A].[Id]) AS [ApplicationId],
		[A].[Name] AS [Application],
		dbo.UidToString([AM].Id) AS [ApplicationModuleId],
		[AM].[Name] AS [Module],
		[I].[RespondedBy],
		[I].[Subject],
		[I].[Body],
		[I].[DateSent],
		[I].[DateResponded],
		[I].[IsApproved],
		[I].[ApproverRemarks],
		[I].[Created],
		[I].[CreatedBy] AS [RequestedBy],
	    isnull(AllowReassign,'') as AllowReassign,
		[IPDR].[IPTitle]
	  FROM [dbo].[Items] AS [I]
		INNER JOIN [dbo].[ApplicationModules] AS [AM] ON [I].[ApplicationModuleId] = [AM].[Id]
		INNER JOIN [dbo].[Applications] AS [A] ON [AM].[ApplicationId] = [A].[Id]
		INNER JOIN [dbo].[IPDisclosureRequest] AS [IPDR] ON [I].[Id] = [IPDR].[ApprovalRequestId]
		INNER JOIN STRING_SPLIT(@Search, ' ') AS ss ON (i.Subject LIKE '%'+ss.value+'%' OR i.CreatedBy LIKE '%'+ss.value+'%')
	  WHERE
		(
			(@IsApproved = 0 AND [I].[IsApproved] IS NULL) OR -- Pending
			(@IsApproved = 1 AND [I].[IsApproved] = 1) OR -- Approved
			(@IsApproved = 2 AND [I].[IsApproved] = 0) OR -- Rejected
			(@IsApproved = 3 AND [I].[IsApproved] IS NOT NULL) -- Closed (Rejected, Approved)
			-- If the value of IsApproved is 4 then select all
		) AND
		[AM].[Id] = @ModuleId
	ORDER BY [I].[Created] DESC
	OFFSET @Offset ROWS 
	FETCH NEXT @Filter ROWS ONLY
END