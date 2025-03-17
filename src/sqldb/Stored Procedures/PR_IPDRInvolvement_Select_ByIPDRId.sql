CREATE PROCEDURE [dbo].[PR_IPDRInvolvement_Select_ByIPDRId]
    @RequestId [INT]
AS
BEGIN
	SELECT
		[RequestId],
		[InvolvementId],
		[I].[Name] AS [InvolvementName]
	FROM
		[dbo].[IPDRInvolvement]
	INNER JOIN
		[dbo].[Involvement] AS [I] ON [I].[Id] = [IPDRInvolvement].[InvolvementId]
	WHERE
		[RequestId] = @RequestId
END

