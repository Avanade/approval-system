CREATE PROCEDURE [dbo].[PR_ApprovalRequestApprovers_Select_ByItemId] 
(
	@ItemId UNIQUEIDENTIFIER
)
AS
BEGIN
	SELECT * FROM [dbo].[ApprovalRequestApprovers] WHERE ItemId = @ItemId
END