CREATE PROCEDURE [dbo].[PR_IPDRInvolvement_Insert]
    @RequestId [INT],
    @InvolvementId [INT]
AS
	INSERT INTO [dbo].[IPDRInvolvement] (
		[RequestId],
        [InvolvementId]
		)
	VALUES (
		@RequestId,
        @InvolvementId
	)

