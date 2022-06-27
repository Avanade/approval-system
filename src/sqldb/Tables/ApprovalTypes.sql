﻿CREATE TABLE [dbo].[ApprovalTypes]
(
	[Id] INT NOT NULL PRIMARY KEY IDENTITY, 
    [Name] VARCHAR(50) NOT NULL,
    [ApproverUserPrincipalName] VARCHAR(100) NULL,
    [IsActive] BIT NOT NULL DEFAULT 1,
    [Created] DATETIME NOT NULL DEFAULT getdate(), 
    [CreatedBy] VARCHAR(100) NULL, 
    [Modified] DATETIME NOT NULL DEFAULT getdate(), 
    [ModifiedBy] VARCHAR(100) NULL
    CONSTRAINT FK_ApprovalTypes_Users FOREIGN KEY (ApproverUserPrincipalName) REFERENCES Users(UserPrincipalName)
)
