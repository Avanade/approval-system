
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
create PROCEDURE [dbo].[PR_Items_Update_ApproverEmail]
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