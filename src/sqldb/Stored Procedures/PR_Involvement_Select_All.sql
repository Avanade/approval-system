CREATE PROCEDURE [dbo].[PR_Involvement_Select_All]
AS
BEGIN
    SELECT
        CONVERT(VARCHAR(36),Id) AS [Id],
        [Name]
    FROM 
        [dbo].[Involvement]
END