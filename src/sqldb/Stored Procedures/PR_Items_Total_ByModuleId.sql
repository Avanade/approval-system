CREATE PROCEDURE [dbo].[PR_Items_Total_ByModuleId]
(
	@Search VARCHAR(50) = '',
	@IsApproved int = 4,
	@ModuleId UNIQUEIDENTIFIER
)
AS
BEGIN
	SELECT
		COUNT([I].[Id]) AS [Total]
	  FROM [dbo].[Items] AS [I]
		INNER JOIN [dbo].[ApplicationModules] AS [AM] ON [I].[ApplicationModuleId] = [AM].[Id]
		INNER JOIN IPDisclosureRequest AS IP ON I.Id = IP.ApprovalRequestId
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
END