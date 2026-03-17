# Project Name: Organization Membership Management Web Application

- **Prepared For:** Google Developer Groups on Campus
- **Prepared by:** Hanz Chrisrome C. Chico
- **Date:** 2026-03-12

## 1. Introduction

**1.1 Purpose:** The purpose of this web application is to enable school organization executives to manage members efficiently, track membership payments. It will streamline administrative tasks, reduce errors from manual record-keeping, and provide a centralized platform for members and for executives to handle membership payments.

**1.2 Scope:**

- The system will allow executives to:
  - Manage member records (CRUD operations)
  - Track payments manually (mark as paid/unpaid)
  - Validate member eligibility for events based on payment status
  - Generate simple reports on members and events
- The system will allow users to:
  - Register as a member.
  - Upload proof of payment
  - View their membership status.
- The system will not handle automated online payments in the initial version but may integrate with online payment gateways in future versions.

## 2. Business Requirements

| ID        | Business Requirement | Description                                                                                                                              |
| :-------- | :------------------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| **BR-01** | Member Registration  | Users must be able to register as a member through an online form, and this registration will let them have access to members dashboard. |
| **BR-02** | Member Management    | Executives must be able to add, view, update, and delete members.                                                                        |
| **BR-03** | Payment Tracking     | Executives must be able to mark membership payments as paid manually.                                                                    |
| **BR-03** | Event Management     | Executives must be able to create events, manage RSVPs, and track attendance.                                                            |
| **BR-05** | Reporting            | The system must generate basic reports for member payments and event attendance.                                                         |
| **BR-06** | Security             | The system must restrict access to executives through secure login.                                                                      |
| **BR-07** | Search & Filter      | Executives should be able to search/filter members and events based on various criteria.                                                 |

## 3. Functional Requirements

- **FR-01: User Authentication**
  - **Description:** Executives and members must log in securely to access the system.
  - **Technical Implementation:** Authentication & role-based access; store hashed passwords; optional JWT or session tokens; rate-limit login attempts; input validation; audit log for login/logout; prevent brute-force attacks.
- **FR-02: Member Registration & Creation**
  - **Description:** Executives can add member details including name, contact info, and membership ID. Users must also be able to register as members and upload proof of payments.
  - **Technical Implementation:** Input validation; uniqueness checks for email/student ID; file upload validation (jpg/png/pdf, max size 5MB); store uploads securely (cloud or DB); default payment status = Pending; audit log on creation; API returns status codes with descriptive messages; concurrent submission handling.
- **FR-03: Update Member**
  - **Description:** Executives can update member information.
  - **Technical Implementation:** Input validation; check uniqueness if email/student ID changed; audit log old vs new values; atomic updates in DB; proper error handling for invalid updates; API returns appropriate status codes.
- **FR-04: Delete Member**
  - **Description:** Executives can remove members from the system.
  - **Technical Implementation:** Validate that member has no pending payments before deletion; soft delete vs hard delete; audit log deletion with actor_id and timestamp; API returns confirmation; secure endpoint access.
- **FR-05: View Members**
  - **Description:** Executives can view a list of all members with payment status.
  - **Technical Implementation:** Pagination, search, and filters; prevent exposing sensitive fields (passwords, internal notes); API supports query parameters; caching for large datasets.
- **FR-06: Record Payment**
  - **Description:** Executives can manually mark a member’s payment as paid.
  - **Technical Implementation:** Validate payment state = Pending; DB update to prevent race conditions; audit log entry
- **FR-07: Search & Filter**
  - **Description:** Executives can search members by name, ID, or payment status.
  - **Technical Implementation:** Backend search supports multiple filters; sanitize inputs to prevent SQL injection; pagination; optional caching for frequent queries; API returns filtered, paginated JSON.
- **FR-08: Reporting**
  - **Description:** Generate simple reports for membership payments.
  - **Technical Implementation:** Backend generates CSV/PDF; filterable by date/payment status/member.
- **FR-09: Audit Logs**
  - **Description:** Executives can view system activity history.
  - **Technical Implementation:** Immutable logs; filterable and paginated; include actor_id, role, action, entity_type, entity_id, timestamp, IP; backend REST API

## 4. Non-Functional Requirements

| ID         | Requirement     | Description                                                                                                 |
| :--------- | :-------------- | :---------------------------------------------------------------------------------------------------------- |
| **NFR-01** | Performance     | Must support up multiple members, rate limiting must be applied to prevent multiple requests in the server. |
| **NFR-02** | Usability       | Interface should be intuitive for non-technical users.                                                      |
| **NFR-03** | Security        | Role-based access control, secure login, and data validation.                                               |
| **NFR-04** | Maintainability | System should allow easy updates to member/event structures.                                                |
| **NFR-05** | Scalability     | Future-ready for more users.                                                                                |

## 5. User Roles

- **Executive**
  - Permissions:
    - Manage Members
    - Approve/Reject Payments
- **Members**
  - Permissions:
    - Register account
    - Upload payment proof

## 6. Use Cases

- **UC-01 – User Registration:**
  - **Actors:** Users. **Precondition:** User is not yet registered in the system.
  - **Main Flow:**
    1. Navigate to “Add Member” page
    2. Enter member details (Name, Contact, Student ID)
    3. Submit form → member added to database
  - **Postcondition:** Member record available in the system
- **UC-02 – User Registration:**
  - **Actors:** Users. **Precondition:** User is not yet registered in the system.
  - **Main Flow:**
    1. User opens the registration page.
    2. User fills out the registration form (Name, Student ID, Email, Password (for account creation), Contact Number, etc.).
    3. User submits the form.
    4. Users attach a proof of payment (screenshot)
    5. System stores the registration information and proof of payment image.
    6. System sets member status to Pending Payment.
  - **Postcondition:** Member record available in the system and the executives can view the attach payments. Users who register must be added to the logs.
- **UC-03 – Mark Payment as Paid:**
  - **Actors:** Executives. **Precondition:** Member exist in the system, and the executive is logged in
  - **Steps:**
    1. Select Member
    2. Mark as paid if the executive handles the registration of the user.
    3. View the attach payment and confirm if the user is really paid.
    4. (Optional) Select a month on the validity of that payment or membership.
  - **Postcondition:** Member is approved in the system, and the member can view in their dashboard that they are successfully a member. Marking a payment must also be added to the logs.
- **UC-04 – Manage Members**
  - **Actor:** Executive. **Preconditions:** Executive is logged in.
  - **Main Flow:**
    1. Executive navigates to the Member Management page.
    2. Executive views the list of members.
    3. Executive can: Add new members, Update member details, Delete members
    4. System saves the changes.
  - **Postconditions:** Member database is updated, and adding a member must be added to the logs with the ID of the added user, and the executive who added the member.
- **UC-05 – View Payment Status**
  - **Actor:** Users. **Preconditions:** User is logged in.
  - **Main Flow:**
    1. User logged in through the web application.
    2. User can view the dashboard with their payment status.
    3. If there is a problem with their proof or they are rejected, user must be aware of the error or mistake and re-submit proof of payment.
    4. The system must change the status to resubmitted.
  - **Postconditions:** User is logged in, and their payment status will be updated, changing payment status must also be added to the audit logs.
- **UC-06 – Search Members**
  - **Actor:** Executive. **Preconditions:** Executive is logged in.
  - **Main Flow:**
    1. Executive opens the member list.
    2. Executive enters search criteria (Name, Student ID, Payment Status).
    3. System displays filtered results.
  - **Postconditions:** Matching members are displayed.
- **UC-07 – Generate Reports**
  - **Actor:** Executive. **Preconditions:** Executive is logged in.
  - **Main Flow:**
    1. Executive navigates to the Reports section.
    2. Executive selects report type (Membership, Payments).
    3. System generates the report.
    4. Executive views or downloads the report.
  - **Postconditions:** Report is generated successfully.

## 7. System Entities

- **Members:** member_id, name, email, student_id, course, contact_number, registration_status, created_at, last_updated
- **Executives:** executive_id, name, email, student_id, course, contact_number, role_id (FK -> Roles.role_id), created_at, last_updated
- **Roles:** role_id, role_name, description, created_at
- **Permissions:** permission_id (PK), permission_key, resource, action, description, created_at, created_by
- **RolePermissions:** role_id (FK → Roles.role_id), permission_id (FK → Permissions.permission_id)
- **Payments:** payment_id, member_id (FK → Members.member_id), payment_proof_image, payment_status, submission_date, approval_date, approved_by (FK → Executives.executive_id)
- **Audit Log:** audit_id, actor_id, actor_role, action, entity_type, entity_id, details, timestamp
- **Sessions:** session_id, owner_id (FK -> Members.member_id OR Executives.executive_id), owner_type (member/executive), refresh_token_hash, user_agent, ip_address, expires_at, created_at, revoked_at


## 8. Technical Design

- **Authentication:**
  - POST /api/login – Executive login
  - POST /api/logout – Logout
- **Members:**
  - POST /api/members – Add member (Executive only)
  - GET /api/members – Get member list with latest payment summary (Executive only)
  - GET /api/members/{id} – Get member details (Executive only)
  - PUT /api/members/{id} – Update member details (Executive only)
  - DELETE /api/members/{id} – Soft delete (Inactivate) member (Executive only)
- **Registration:**
  - POST /api/register – User registration
  - POST /api/payments – Upload payment proof
  - GET /api/payments/{id} – Get payment submission
  - PUT /api/payments/{id}/approve – Approve payment
  - PUT /api/payments/{id}/reject – Reject payment

## 9. Frontend / UI Design

- **10.1 Login page (Executives)**
  - **Components:** Email / Username input, Password input, Login button, “Forgot Password?” link, Validation messages (invalid credentials, empty fields), Optional: Show login error notifications (toast or inline)
  - **Behavior / Notes:** On successful login → redirect to Executive Dashboard, On failure → display error, Mobile-friendly responsive form
- **10.2 Member Registration Page (Self-Registration)**
  - **Components:** Form fields: Full Name, Email Address, Student ID, Course / Program, Contact Number. Upload section for payment screenshot, Submit button, Validation messages (empty fields, invalid email), Success confirmation message.
  - **Behavior / Notes:** Default membership status: Pending Payment, File type validation (JPG, PNG), Maximum file size limit
- **10.3 Executive Dashboard**
  - **Components / Sections:**
    - Sidebar Navigation: Members, Payment Approvals, Events, Reports, Logout
    - Quick Stats / Widgets: Total Members, Pending Payments, Upcoming Events
    - Tables / Lists: Member list table with columns: Name, Student ID, Email, Payment Status, Actions (Edit/Delete/View); Payment approval table with screenshot preview and Approve / Reject buttons
    - Search & Filter Components: By Name, Student ID, Payment Status, Event Name
  - **Behavior / Notes:** Clickable actions → open modals for editing / viewing, Table supports pagination and sorting, Responsive design for tablets / small screens
- **10.4 Member Management Page (Executives)**
  - **Components:** Member list table (as above), Add New Member button → opens modal or page with form: Name, Email, Student ID, Contact, Course. Edit button → opens modal pre-filled with existing data, Delete button → confirmation dialog
  - **Behavior / Notes:** Real-time updates after Add / Edit / Delete, Payment status clearly displayed for each member
- **10.5 Payment Approval Page (Executives)**
  - **Components:** Table of submitted payments: Member Name, Student ID, Payment Proof (thumbnail or click-to-view), Status (Pending / Approved / Rejected). Approve / Reject buttons, Search & Filter: by Name, Status, Pagination if many submissions, Notification / toast after approve/reject
  - **Behavior / Notes:** Upon approval, delete the image attach in the database to prevent exceeding storage limits. Rejection allows member to upload a new proof.
- **10.6 Members Dashboard Page (Users)**
  - **Components:** Simple page only that shows the payment status, and their details.
- **9. Reports Page (Executives)**
  - **Components:** Dropdown to select report type: Members list, Payment status. Filter options: By date, payment status, usernames, student ID. Generate / Download button, Table or chart visualization, Pagination for large data sets
  - **Behavior / Notes:** Allow CSV / PDF download, Optional graphs for dashboard summary
- **Audit logs page (Executives):**
  - **Components:**
    - 1. Filters / Controls
    - Dropdown to select log type / category: All Actions, Member Actions, Executive Actions, System Actions (optional)
    - Search & Filter Options: By actor name or username, By role (Member / Executive), By action type (e.g., Login, Add Member, Approve Payment, RSVP), By entity type (Member, Payment, Event, RSVP), By date range (start date / end date)
    - Generate / Refresh button to apply filters
    - Table / Lists: Timestamp, actor name, role, action, entity type, entity ID, details.
    - Pagination support, Sortable columns (sort by timestamp)
  - **Behavior / Notes:** Logs display in reverse chronological order (most recent first), Ensure only executives/admins can access audit logs, Table should be searchable and filterable in real-time

## 10. UX / UI Notes

- Add a header always for each page, (e.g., Welcome back, Executive, View Members)
- Clean card-based layout for quick info
- Navigation: Sidebar for executives, top nav for members
- Notifications: Toast messages for success/error
- Modals: Use modals for Add/Edit forms to keep dashboard clean
- Status indicators use consistent color codes: Pending → Yellow, Approved → Green, Rejected → Red
- Responsive design for mobile and tablet
- Highlight actions that require member input (e.g., Upload Payment, RSVP)
- Theme must be based on: Google Developer Groups on Campus.

## 11. System Architecture

- **Frontend:** React, Typescript
- **Frontend Framework:** TailwindCSS
- **Backend:** Go
- **Database:** Supabase
- **Authentication:** JWT or Session-based login
- **Hosting:** Vercel (Frontend) Render (Backend)
