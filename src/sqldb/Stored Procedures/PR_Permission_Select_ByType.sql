CREATE PROCEDURE  [dbo].[PR_Permission_Select_ByType]
(
    @Type VARCHAR(100)
)
AS
BEGIN   
    SELECT 
        [P].[Email]
    FROM [dbo].[Permission] AS [P]
    WHERE
        [P].[Type] = @Type
END