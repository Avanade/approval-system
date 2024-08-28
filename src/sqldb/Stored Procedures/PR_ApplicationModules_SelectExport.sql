CREATE PROCEDURE PR_ApplicationModules_SelectExport_ById
AS
BEGIN
    SELECT
        ExportUrl
    FROM
        ApplicationModules
    WHERE
        ExportUrl IS NOT NULL
END