CREATE PROCEDURE [dbo].[PR_Items_UpdateRespondedBy_ById]
(
    @Id UNIQUEIDENTIFIER,
    @RespondedBy VARCHAR(100)
)
AS
BEGIN
    -- SET NOCOUNT ON added to prevent extra result sets from
    -- interfering with SELECT statements.
    SET NOCOUNT ON

    -- Insert statements for procedure here
    UPDATE [dbo].[Items]
    SET [RespondedBy] = @RespondedBy,
        [Modified] = GETDATE(),
        [ModifiedBy] =  @RespondedBy
    WHERE [Id] = @Id
END