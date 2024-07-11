CREATE PROCEDURE [dbo].[PR_Items_Update_ApproverEmail]
	@Id uniqueidentifier,
	@ApproverEmail varchar(100),
 	@Username VARCHAR(100)
AS
BEGIN
	UPDATE 
		[dbo].[ApprovalRequestApprovers]
	SET
		ApproverEmail = @ApproverEmail
	WHERE 
		ItemId = @Id AND ApproverEmail = @Username;

	UPDATE 
		[dbo].[Items]
	SET
 		Modified = GETDATE(),
		ModifiedBy = @Username
	WHERE 
		Id = @Id;
END