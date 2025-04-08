CREATE PROCEDURE [dbo].[PR_IPDisclosureRequest_Insert]
    @RequestorName [VARCHAR](100),
    @RequestorEmail [VARCHAR](100),
    @IPTitle [VARCHAR](100),
    @IPType [VARCHAR](100),
    @IPDescription [VARCHAR](350),
    @Reason [VARCHAR](350)
AS

	DECLARE @ResultTable TABLE(Id [INT]);

	INSERT INTO [dbo].[IPDisclosureRequest] (
		[RequestorName],
        [RequestorEmail],
        [IPTitle],
        [IPType],
        [IPDescription],
        [Reason],
        [Created]
		)
	OUTPUT INSERTED.Id INTO @ResultTable
	VALUES (
		@RequestorName,
        @RequestorEmail,
        @IPTitle,
        @IPType,
        @IPDescription,
        @Reason,
        GETDATE()
	)

	SELECT [Id] FROM @ResultTable

