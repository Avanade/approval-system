CREATE TABLE [dbo].[ApprovalRequestApprovers]
(   
    [ItemId] UNIQUEIDENTIFIER NOT NULL,
    [ApproverEmail] VARCHAR(100) NOT NULL
    CONSTRAINT PK_ApprovalRequestApprover PRIMARY KEY (ItemId, ApproverEmail),
    CONSTRAINT FK_ApprovalRequestApprover_Items FOREIGN KEY (ItemId) REFERENCES Item(Id)
)