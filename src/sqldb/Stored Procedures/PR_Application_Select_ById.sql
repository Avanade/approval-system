CREATE PROCEDURE PR_Applications_Select_ById

	@Id UNIQUEIDENTIFIER

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