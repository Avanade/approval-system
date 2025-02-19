CREATE PROCEDURE  [dbo].[PR_Permission_Insert]
(
    @Type VARCHAR(100),
    @Email VARCHAR(100),
    @CreatedBy VARCHAR(100)
)
AS
BEGIN   
   INSERT INTO [dbo].[Permission] 
    (
        [Type], 
        [Email], 
        [Created], 
        [CreatedBy]
    )
    VALUES 
    (
        @Type, 
        @Email, 
        GETDATE(), 
        @CreatedBy
    )
END