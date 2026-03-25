# FoodStore

FoodStore is a marketplace-style platform for food commerce where customers, sellers, and platform teams interact in one ecosystem.

The idea is simple: help people buy and sell food products in a clear, convenient, and trustworthy environment.

## What This Project Is About

FoodStore is focused on practical user actions:
- A customer can browse available products and place orders.
- A seller can add their own products, manage listings, and grow their store visibility.
- A moderator can review marketplace content and support quality standards.
- An admin can manage platform-level operations and user trust.

## Why FoodStore Is Strong

- **Role-oriented experience**: each user type gets the right tools for their goals.
- **Clear marketplace journey**: product discovery, ordering, and management are straightforward.
- **Trust-focused structure**: moderation and administration help keep the platform reliable.
- **Growth potential**: designed to support both new and expanding food sellers.

## Tech Stack

- **Language**: Go
- **API**: HTTP REST
- **Routing/Delivery Layer**: custom HTTP delivery with handlers and middleware
- **Database**: PostgreSQL
- **Configuration**: environment variables (`envconfig`, `.env`)
- **Authentication**: JWT-based auth flow
- **Security**: password hashing and code generation utilities
- **Email**: SMTP sender for recovery and notifications
- **Architecture Style**: layered/clean-style separation (`delivery`, `usecase`, `domain`, `infrastructure`)
- **DevOps**: Docker, Docker Compose, Makefile

## Product Vision

FoodStore aims to become a practical and user-friendly place where independent food sellers can present their products, and customers can quickly find and order what they need.

## Current Status

FoodStore is still in development and not fully finished yet - we are waiting for the official release.
