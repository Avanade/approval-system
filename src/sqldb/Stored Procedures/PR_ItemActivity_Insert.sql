CREATE PROCEDURE [dbo].[PR_ItemActivity_Insert]
    @CreatedBy VARCHAR(50),
    @Content VARCHAR(8000),
    @ItemId UNIQUEIDENTIFIER
AS
BEGIN
    DECLARE @ResultTable TABLE(Id [INT]);

	INSERT INTO [dbo].[ItemActivity] (
		[CreatedBy],
        [Content],
        [ItemId],
        [Created]
		)
    OUTPUT INSERTED.Id 
	VALUES (
		@CreatedBy,
        @Content,
        @ItemId,
        GETDATE()
	)

    SELECT [Id] FROM @ResultTable
END

