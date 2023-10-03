CREATE PROCEDURE [dbo].[PR_Items_Select_All]
AS
BEGIN
    SELECT
        CONVERT(VARCHAR(36),Id) AS Id,
        ApproverEmail,
        DateResponded
    FROM 
        [dbo].[Items]
END