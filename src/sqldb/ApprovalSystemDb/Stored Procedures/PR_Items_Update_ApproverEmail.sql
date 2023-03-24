

CREATE PROCEDURE [dbo].[PR_Items_Update_ApproverEmail]
	@Id uniqueidentifier,
	@ApproverEmail varchar(100),
 	@Username VARCHAR(100)
AS
	UPDATE Items
	SET
		ApproverEmail = @ApproverEmail,
 		Modified = GETDATE(),
		ModifiedBy = @Username
	WHERE Id = @Id