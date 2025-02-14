CREATE PROCEDURE  [dbo].[PR_LegalConsultation_Select_ByEmail]
(
    @Email VARCHAR(100),
    @Offset INT = 0,
    @Filter INT = 10,
    @IsApproved INT = 4
)
AS
BEGIN   
    SELECT 
        [dbo].UidToString([LC].[ItemId]) AS [ItemId],
        [LC].[Created] AS [LegalConsultationCreated],
        [dbo].UidToString([A].[Id]) AS [ApplicationId],
        [A].[Name] AS [Application],
        [dbo].UidToString([AM].[Id]) AS [ApplicationModuleId],
        [AM].[Name] AS [Module],
        [I].[Subject],
        [I].[Body],
        [I].[DateSent],
        [I].[DateResponded],
        [I].[IsApproved],
        [I].[ApproverRemarks],
        [I].[RespondedBy],
        [I].[Created],
        [I].[CreatedBy]
    FROM [dbo].[LegalConsultation] AS [LC]
    INNER JOIN [dbo].[Items] AS [I] ON [I].[Id] = [LC].[ItemId]
    INNER JOIN [dbo].[ApplicationModules] AS [AM] ON [I].[ApplicationModuleId] = [AM].[Id]
    INNER JOIN [dbo].[Applications] AS [A] ON [AM].[ApplicationId] = [A].[Id]
    WHERE [LC].[Email] = @Email AND
		(
			(@IsApproved = 0 AND i.IsApproved IS NULL) OR -- Pending
			(@IsApproved = 1 AND i.IsApproved = 1) OR -- Approved
			(@IsApproved = 2 AND i.IsApproved = 0) OR -- Rejected
			(@IsApproved = 3 AND i.IsApproved IS NOT NULL) -- Closed (Rejected, Approved)
			-- If the value of IsApproved is 4 then select all
		)
    ORDER BY [I].[Created] DESC
    OFFSET @Offset ROWS 
    FETCH NEXT @Filter ROWS ONLY
END