CREATE PROCEDURE [dbo].[PR_ApplicationModules_Select]
AS
BEGIN
    SET NOCOUNT ON
    SELECT
        CONVERT(varchar(36), [Id], 1) AS [Id],
        [Name],
		[CallbackUrl], 
		[RequireRemarks], 
		[ReassignCallbackUrl],
		[ExportUrl],
		[AllowReassign],
		[RequireAuthentication]
    FROM [dbo].[ApplicationModules]
END