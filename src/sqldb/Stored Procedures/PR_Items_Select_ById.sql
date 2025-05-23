CREATE PROCEDURE [dbo].[PR_Items_Select_ById]
	@Id UNIQUEIDENTIFIER
AS
	SELECT
		A.[Name] [Application],
		dbo.UidToString(A.[Id]) AS [ApplicationId],
		AM.[Name] [Module], dbo.UidToString(AM.[Id]) [ApplicationModuleId],
		[Subject], Body, DateSent,
		DateResponded, IsApproved, ApproverRemarks,
		T.ApproveText, T.RejectText,
		AM.CallbackUrl,
	    AM.ReassignCallbackUrl,
		i.RespondedBy,
		I.Created, I.CreatedBy [RequestedBy]
	FROM Items I
	INNER JOIN ApplicationModules AM ON I.ApplicationModuleId = AM.Id
	INNER JOIN Applications A ON AM.ApplicationId = A.Id
	INNER JOIN ApprovalTypes T ON AM.ApprovalTypeId = T.Id
	WHERE I.Id = @Id
	ORDER BY I.Created DESC