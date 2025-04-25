CREATE TABLE [dbo].[IPDisclosureRequest]
(   
    [Id]  INT          IDENTITY (1, 1),
    [RequestorName] VARCHAR(100) NOT NULL,
    [RequestorEmail] VARCHAR(100) NOT NULL,
    [IPTitle] VARCHAR(100) NOT NULL,
    [IPType] VARCHAR(100) NOT NULL,
    [IPDescription] VARCHAR(1000) NOT NULL,
    [Reason] VARCHAR(1000) NOT NULL,
    [ApprovalRequestId] UNIQUEIDENTIFIER NULL,
    [IsApproved] BIT NULL,
    [ApproverRemarks] VARCHAR(255) NULL,
    [Created] DATETIME NOT NULL,
    [ResponseDate] DATETIME NULL,
    [RespondedBy] VARCHAR(100) NULL
    CONSTRAINT PK_IPDisclosureRequest PRIMARY KEY (Id)
)