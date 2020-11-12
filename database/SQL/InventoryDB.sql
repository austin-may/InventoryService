DROP DATABASE IF EXISTS [InventoryDB]
GO
CREATE DATABASE [InventoryDB]

GO

USE [InventoryDB]

GO

CREATE TABLE InventoryReserve (
    [ItemID] INT IDENTITY (1,1),
    [Name] VARCHAR(50),
    [Count] INT,
    [Site] VARCHAR(50)
	PRIMARY KEY ([ItemID])
)