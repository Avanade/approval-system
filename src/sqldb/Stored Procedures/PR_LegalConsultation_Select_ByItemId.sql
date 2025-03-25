CREATE PROCEDURE  [dbo].[PR_LegalConsultation_Select_ByItemId]
(
    @ItemId UNIQUEIDENTIFIER
)
AS
BEGIN   
    SELECT 
        [Email],
        [Created],
        [CreatedBy]
    FROM [dbo].[LegalConsultation]
    WHERE [ItemId] = @ItemId
END