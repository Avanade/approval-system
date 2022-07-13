CREATE PROCEDURE PR_ApplicationModules_IsAuthRequired
    @ApplicationModuleId UNIQUEIDENTIFIER

AS

SELECT RequireAuthentication FROM ApplicationModules WHERE Id = @ApplicationModuleId