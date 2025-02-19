CREATE PROCEDURE  [dbo].[PR_LegalConsultation_Insert]
(
    @ItemId UNIQUEIDENTIFIER,
    @Email VARCHAR(100),
    @CreatedBy VARCHAR(100)
)
AS
BEGIN   
    INSERT INTO [dbo].[LegalConsultation]
        (
            [ItemId],
            [Email],
            [Created],
            [CreatedBy]
        )
    VALUES
        (
            @ItemId,
            @Email,
            GETDATE(),
            @CreatedBy
        )
END