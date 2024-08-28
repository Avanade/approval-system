CREATE PROCEDURE  [dbo].[PR_ApprovalRequestApprovers_Insert]
(
    @ItemId UNIQUEIDENTIFIER,
    @ApproverEmail VARCHAR(100)
)
AS
BEGIN   
    INSERT INTO [dbo].[ApprovalRequestApprovers]
        (
            [ItemId],
            [ApproverEmail]
        )
    VALUES
        (
            @ItemId,
            @ApproverEmail
        )
END
