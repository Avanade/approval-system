-- This file contains SQL statements that will be executed after the build script.
/* INITIAL DATA FOR APPROVAL TYPES */
SET IDENTITY_INSERT ApprovalTypes ON

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (1, 'Approve', 'Reject')

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 2)
INSERT INTO ApprovalTypes (Id, ApproveText, RejectText) VALUES (2, 'Accept', 'Decline')

SET IDENTITY_INSERT ApprovalTypes OFF

SET IDENTITY_INSERT Involvement ON

IF NOT EXISTS (SELECT Id FROM Involvement WHERE Id = 1)
INSERT INTO Involvement (Id, Name) VALUES (1, 'MVP')

IF NOT EXISTS (SELECT Id FROM Involvement WHERE Id = 2)
INSERT INTO Involvement (Id, Name) VALUES (2, 'RD')

IF NOT EXISTS (SELECT Id FROM Involvement WHERE Id = 3)
INSERT INTO Involvement (Id, Name) VALUES (3, 'Open Source')

IF NOT EXISTS (SELECT Id FROM Involvement WHERE Id = 4)
INSERT INTO Involvement (Id, Name) VALUES (4, 'External Speaker')

SET IDENTITY_INSERT Involvement OFF