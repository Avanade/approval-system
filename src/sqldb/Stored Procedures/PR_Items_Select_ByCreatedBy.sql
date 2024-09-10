CREATE PROCEDURE [dbo].[PR_Items_Select_ByCreatedBy]
	@CreatedBy varchar(255)
AS
	SELECT
		dbo.UidToString(A.Id) [ApplicationId], A.[Name] [Application],
		AM.[Name] [Module],
		[Subject], Body, DateSent,
		DateResponded, IsApproved, ApproverRemarks, I.CreatedBy, I.Created
	FROM Items I
	INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
	INNER JOIN Applications A ON AM.ApplicationId = A.Id
	WHERE I.CreatedBy = @CreatedBy
	ORDER BY I.Created DESC