CREATE PROCEDURE [dbo].[PR_IPDisclosureRequest_Update_ApprovalRequestId]
  @ApprovalRequestId [UNIQUEIDENTIFIER],
  @IPDRequestId [INT]
AS
BEGIN
  UPDATE [dbo].[IPDisclosureRequest]
  SET
    [ApprovalRequestId] = @ApprovalRequestId
  WHERE [Id] = @IPDRequestId
END
