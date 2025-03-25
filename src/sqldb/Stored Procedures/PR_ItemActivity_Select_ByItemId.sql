CREATE PROCEDURE [dbo].[PR_ItemActivity_Select_ByItemId]
    @ItemId [UNIQUEIDENTIFIER]
AS
BEGIN
	SELECT
        [Id],
        [CreatedBy],
        [Created],
        [Content]
    FROM  [dbo].[ItemActivity]
    WHERE [ItemId] = @ItemId
END

