-- This file contains SQL statements that will be executed after the build script.

/* INITIAL DATA FOR APPROVAL STATUS */

SET IDENTITY_INSERT ApprovalStatus ON

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 1)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (1, 'New')

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 2)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (2, 'InReview')

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 3)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (3, 'Rejected')

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 4)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (4, 'NonCompliant')

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 5)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (5, 'Approved')

IF NOT EXISTS (SELECT Id FROM ApprovalStatus WHERE Id = 6)
INSERT INTO ApprovalStatus (Id, [Name]) VALUES (6, 'Retired')

SET IDENTITY_INSERT ApprovalTypes OFF

/* INITIAL DATA FOR APPROVAL TYPES */

SET IDENTITY_INSERT ApprovalTypes ON

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 1)
INSERT INTO ApprovalTypes (Id, [Name]) VALUES (1, 'Intellectual Property')

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 2)
INSERT INTO ApprovalTypes (Id, [Name]) VALUES (2, 'Legal')

IF NOT EXISTS (SELECT Id FROM ApprovalTypes WHERE Id = 3)
INSERT INTO ApprovalTypes (Id, [Name]) VALUES (3, 'Security')

SET IDENTITY_INSERT ApprovalTypes OFF