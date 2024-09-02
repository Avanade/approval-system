CREATE PROCEDURE PR_Applications_Select_ById

	@Id UNIQUEIDENTIFIER = 'ba8bfbd0-14c9-45da-af6d-2d0839924806'

AS
BEGIN
	SELECT
		dbo.UidToString(Id) [Id],
		[Name],
		[ExportUrl],
		[OrganizationTypeUrl]
	FROM
		Applications
	WHERE
		IsActive = 1 AND Id = @Id
END