CREATE PROCEDURE PR_ApplicationModules_Select_ById_ApplicationId(
	@ApplicationId UNIQUEIDENTIFIER,
	@ApplicationModuleId UNIQUEIDENTIFIER
)
AS
BEGIN

	SELECT 
		[A].[Name] AS [ApplicationName], 
		[AM].[Name] AS [ApplicationModuleName],
		[AM].[CallbackUrl], 
		[AM].[RequireRemarks], 
		[AT].[ApproveText], 
		[AT].[RejectText],
		[AM].[ReassignCallbackUrl],
		[AM].[ExportUrl],
		[AM].[AllowReassign],
		[AM].[RequireAuthentication]
	FROM [dbo].[ApplicationModules] AS [AM]
	INNER JOIN [dbo].[Applications] AS [A] ON [A].[Id] = [AM].[ApplicationId]
	INNER JOIN [dbo].[ApprovalTypes] AS [AT] ON [AT].[Id] = [AM].[ApprovalTypeId]
	WHERE [A].[IsActive] = 1
		AND [AM].[IsActive] = 1
		AND [A].[Id] = @ApplicationId
		AND [AM].[Id] = @ApplicationModuleId
END