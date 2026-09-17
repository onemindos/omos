# Zoho WorkDrive

> **Cloud Storage & Team Collaboration**

## Overview

Zoho WorkDrive is a cloud file management and collaboration platform for storing, organizing, sharing, and working on files together with your team.

## Status

- **Active**: ✅ Yes
- **Integration**: API + Desktop/Mobile Apps
- **Last Updated**: January 25, 2026

## Use Cases

1. **File Storage** - Central repository for all files
2. **Team Collaboration** - Shared folders and files
3. **Version Control** - Track file changes
4. **External Sharing** - Share files with clients
5. **Real-Time Editing** - Collaborate on documents

## Key Features Used

- [ ] Team folders (shared spaces)
- [ ] File versioning
- [ ] External sharing links
- [ ] Desktop sync client
- [ ] Mobile apps (iOS/Android)
- [ ] Online document editing (Writer, Sheet, Show)
- [ ] File previews (50+ formats)
- [ ] Search & tagging

## Storage Structure

```
OneMind WorkDrive/
├── 01-Projects/
│   ├── Project-Alpha/
│   ├── Project-Beta/
│   └── Client-Work/
├── 02-Documents/
│   ├── Proposals/
│   ├── Contracts/
│   └── Invoices/
├── 03-Marketing/
│   ├── Brand-Assets/
│   ├── Campaigns/
│   └── Content/
├── 04-Internal/
│   ├── SOPs/
│   ├── Templates/
│   └── Training/
└── 05-Archive/
```

## Storage Capacity

| Plan             | Storage | Active Users |
| ---------------- | ------- | ------------ |
| **Free**         | 5 GB    | 25 users     |
| **Standard**     | 1 TB    | Unlimited    |
| **Professional** | 3 TB    | Unlimited    |
| **Current**      | 1 TB    | 5 users      |

## API Integration

### Upload File

```python
POST https://workdrive.zoho.com/api/v1/upload
Headers:
  Authorization: Zoho-oauthtoken {access_token}
Body:
  file: <binary>
  parent_id: {folder_id}
```

### Download File

```python
GET https://workdrive.zoho.com/api/v1/download/{file_id}
```

### Create Folder

```python
POST https://workdrive.zoho.com/api/v1/folders
{
  "name": "New Folder",
  "parent_id": "{parent_folder_id}"
}
```

### List Files

```python
GET https://workdrive.zoho.com/api/v1/files?parent_id={folder_id}
```

## File Sharing

### Internal Sharing

- Share with team members by email
- Set permissions: View, Download, Edit, Upload
- Access levels: Owner, Admin, Member, Viewer

### External Sharing

- Generate share links
- Password protection
- Expiration dates
- Download limits
- Email verification

### Example Share Settings

- **Client Deliverables**: View + Download, 30-day expiry
- **Feedback Request**: View + Comment only
- **Public Resource**: Anyone with link, no expiry

## Version Control

### Automatic Versioning

- All file edits create new versions
- Keep last 50 versions (configurable)
- Restore previous versions
- Compare versions side-by-side
- Download specific version

### Version Naming

- v1.0, v1.1, v1.2... (automatic)
- Add version comments
- Tag major versions

## Desktop Sync Client

### Features

- Two-way sync (desktop ↔ cloud)
- Selective sync (choose folders)
- Offline access
- Background upload/download
- Conflict resolution

### Sync Settings

```
Sync Folder: ~/OneMind-WorkDrive/
Sync Frequency: Real-time
Bandwidth Limit: Unlimited
Exclude: .DS_Store, Thumbs.db, *.tmp
```

## Mobile Apps

### Mobile Features

- View all files
- Upload from camera/gallery
- Offline access (marked files)
- Share files
- Comment on files
- Scan documents (camera → PDF)

### Mobile Optimization

- Auto-upload photos
- Cell data usage limits
- Wi-Fi only uploads

## Online Editing

### Supported Formats

- **Zoho Writer**: Documents (.docx, .odt, .rtf)
- **Zoho Sheet**: Spreadsheets (.xlsx, .ods, .csv)
- **Zoho Show**: Presentations (.pptx, .odp)

### Collaboration Features

- Real-time co-editing
- Track changes
- Comments & mentions
- Version history
- Auto-save (every 30 seconds)

## File Previews

### Supported Formats (50+)

- **Documents**: PDF, DOCX, TXT, MD
- **Spreadsheets**: XLSX, CSV
- **Presentations**: PPTX
- **Images**: JPG, PNG, GIF, SVG
- **Videos**: MP4, MOV, AVI
- **Audio**: MP3, WAV
- **Code**: PY, JS, HTML, CSS, JSON
- **Archives**: ZIP, RAR

## Search & Organization

### Search Features

- Full-text search
- File name/content search
- Filter by: Type, Owner, Date
- Advanced search operators

### Tagging System

- Custom tags for files/folders
- Color-coded labels
- Tag-based filtering
- Tag suggestions

## Integrations

- **Zoho CRM**: Attach files to records
- **Zoho Desk**: Attach to tickets
- **Zoho Mail**: Send as attachment
- **Zoho Projects**: Link project files
- **Google Drive Import**: Migrate files
- **Dropbox Import**: Migrate files

## Team Folders (Team Spaces)

### Structure

Each team folder has:

- **Members**: Users with access
- **Roles**: Owner, Admin, Member, Viewer
- **Sub-folders**: Organize by project/topic
- **Activity Log**: Who did what

### Permissions

- **Owner**: Full control, delete team folder
- **Admin**: Manage members, settings
- **Member**: Upload, edit, delete own files
- **Viewer**: Read-only access

## Automation & Workflows

### Automated Actions

1. **New File Uploaded** → Notify team in Cliq
2. **File Expired** → Move to Archive folder
3. **Large Upload** → Send email notification
4. **External Share Created** → Log in audit trail

## Legacy AI Integration

### Use Cases

1. **Document Processing** - Extract text from uploads
2. **File Organization** - Auto-tag/categorize files
3. **Duplicate Detection** - Find duplicate files
4. **Archive Management** - Move old files to archive
5. **Backup Sync** - Mirror critical files to local storage

## Security Features

- ✅ AES 256-bit encryption (at rest)
- ✅ SSL/TLS encryption (in transit)
- ✅ Two-factor authentication
- ✅ File access audit logs
- ✅ Admin controls (restrict downloads, sharing)
- ✅ Data residency options (US, EU, India)
- ✅ GDPR compliant

## Backup Strategy

### Current Setup

- **Primary**: Zoho WorkDrive (cloud)
- **Secondary**: Local sync (desktop) - Auto-backup
- **Tertiary**: AWS S3 (critical files) - Weekly backup

### Retention Policy

- Active files: Forever
- Deleted files: 60 days in trash
- Old versions: Last 50 versions

## Analytics & Activity

### Tracked Metrics

- Storage usage by folder
- Most active files
- User activity (uploads, downloads, edits)
- External shares count
- Collaboration activity

## Notes

- Free tier: 5 GB storage, 25 users
- Paid plans: $2.50/user/month (Standard), $4/user/month (Professional)
- Desktop sync client available for Windows, Mac, Linux
- Mobile apps on iOS and Android (free)
- Online editing requires Zoho Workplace plan
- API rate limits: 600 calls/minute

---

**Related**: [Cliq](./Cliq.md) | [Mail](./Mail.md)
