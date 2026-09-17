### Create Object Store Bucket

Source: https://github.com/nats-io/nats.docs/blob/master/nats-concepts/jetstream/object-store/obj_walkthrough.md

Initializes a new Object Store bucket named 'myobjbucket' within the JetStream environment.

```shell
nats object add myobjbucket
```

--------------------------------

### ObjectStore API

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Provides methods for performing CRUD operations on objects within an object store, as well as retrieving object store status.

```APIDOC
## ObjectStore API

### Description
The ObjectStore interface contains methods to operate on an object store. It allows for CRUD operations on objects, retrieval of object information, metadata updates, linking objects, watching for updates, listing objects, and retrieving store status and configuration.

### Methods

#### Put

*   **Description**: Places the contents from a reader into a new object. Overwrites the object if it already exists. The object name is taken from `ObjectMeta.Name`.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `obj` (ObjectMeta) - Metadata for the object, including its name.
    *   `reader` (io.Reader) - The data to be stored.
*   **Returns**: `*ObjectInfo` containing the object's metadata, digest, and instance information, and an `error`.

#### PutBytes

*   **Description**: Convenience function to put a byte slice into this object store under the given name.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `name` (string) - The name of the object.
    *   `data` ([]byte) - The byte slice data to store.
*   **Returns**: `*ObjectInfo` containing the object's metadata, digest, and instance information, and an `error`.

#### PutString

*   **Description**: Convenience function to put a string into this object store under the given name.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `name` (string) - The name of the object.
    *   `data` (string) - The string data to store.
*   **Returns**: `*ObjectInfo` containing the object's metadata, digest, and instance information, and an `error`.

#### PutFile

*   **Description**: Convenience function to put the contents of a file into this object store. The object's name will be the path of the file.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `file` (string) - The path to the file.
*   **Returns**: `*ObjectInfo` containing the object's metadata, digest, and instance information, and an `error`.

#### Get

*   **Description**: Retrieves the named object from the object store.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `name` (string) - The name of the object to retrieve.
*   **Returns**: `*ObjectResult` containing the object's metadata and a reader for its contents, and an `error`. Returns `ErrObjectNotFound` if the object does not exist.
```

--------------------------------

### Manage Object Stores with Go

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

The ObjectStoreManager interface provides methods to create, update, delete, and list object store buckets. It is the primary entry point for managing the lifecycle of object stores within a NATS JetStream context.

```go
type ObjectStoreManager interface {
	ObjectStore(ctx context.Context, bucket string) (ObjectStore, error)
	CreateObjectStore(ctx context.Context, cfg ObjectStoreConfig) (ObjectStore, error)
	UpdateObjectStore(ctx context.Context, cfg ObjectStoreConfig) (ObjectStore, error)
	CreateOrUpdateObjectStore(ctx context.Context, cfg ObjectStoreConfig) (ObjectStore, error)
	DeleteObjectStore(ctx context.Context, bucket string) error
	ObjectStoreNames(ctx context.Context) ObjectStoreNamesLister
	ObjectStores(ctx context.Context) ObjectStoresLister
}
```

--------------------------------

### Get Object Store Status - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Retrieves the run-time status of the backing store for the object bucket. This provides insights into the operational state of the store.

```csharp
ValueTask<NatsObjStatus> GetStatusAsync(CancellationToken cancellationToken = default);
```

--------------------------------

### Seal Object Store - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Seals the object store, preventing any further modifications. This operation is useful for making the store immutable after all desired objects have been added. It may throw an NatsObjException if the update operation fails.

```csharp
ValueTask SealAsync(CancellationToken cancellationToken = default);
```

--------------------------------

### C# Object Store Creation and Access

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

C# interface for NATS Object Store operations. Provides methods to create a new object store using either a bucket name or a configuration object, and to retrieve an existing object store.

```C#
// dotnet add package NATS.Net

/// <summary>
/// NATS Object Store context.
/// </summary>
public interface INatsObjContext
{
    /// <summary>
    /// Provides access to the JetStream context associated with the Object Store operations.
    /// </summary>
    INatsJSContext JetStreamContext { get; }

    /// <summary>
    /// Create a new object store.
    /// </summary>
    /// <param name="bucket">Bucket name.</param>
    /// <param name="cancellationToken">A <see cref="CancellationToken"/> used to cancel the API call.</param>
    /// <returns>Object store object.</returns>
    ValueTask<INatsObjStore> CreateObjectStoreAsync(string bucket, CancellationToken cancellationToken = default);

    /// <summary>
    /// Create a new object store.
    /// </summary>
    /// <param name="config">Object store configuration.</param>
    /// <param name="cancellationToken">A <see cref="CancellationToken"/> used to cancel the API call.</param>
    /// <returns>Object store object.</returns>
    ValueTask<INatsObjStore> CreateObjectStoreAsync(NatsObjConfig config, CancellationToken cancellationToken = default);

    /// <summary>
    /// Get an existing object store.
    /// </summary>
    /// <param name="bucket">Bucket name</param>
    /// <param name="cancellationToken">A <see cref="CancellationToken"/> used to cancel the API call.

```

--------------------------------

### Python Object Store Management

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Python functions for managing NATS Object Stores. Includes methods to create and delete object stores, and to retrieve an existing object store by its bucket name.

```Python
async def object_store(self, bucket: str) -> ObjectStore:

    async def create_object_store(
        self,
        bucket: str = None,
        config: Optional[api.ObjectStoreConfig] = None,
        **params,
    ) -> ObjectStore:
        """
        create_object_store takes an api.ObjectStoreConfig and creates a OBJ in JetStream.
        """
    async def delete_object_store(self, bucket: str) -> bool:
        """
        delete_object_store will delete the underlying stream for the named object.
        """
```

--------------------------------

### POST /object-store/create

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Creates a new Object Store bucket in JetStream.

```APIDOC
## POST /object-store/create

### Description
Creates a new object store bucket in JetStream with the specified configuration.

### Method
POST

### Endpoint
/object-store/create

### Parameters
#### Request Body
- **bucket** (string) - Required - The name of the bucket to create.
- **config** (ObjectStoreConfig) - Optional - Configuration settings for the object store.

### Request Example
{
  "bucket": "my-bucket",
  "config": { "description": "My storage bucket" }
}

### Response
#### Success Response (200)
- **ObjectStore** (Object) - The created object store instance.

#### Response Example
{
  "status": "success",
  "bucket": "my-bucket"
}
```

--------------------------------

### List Objects in Store - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Lists all objects currently present in the object store. It returns an asynchronous enumerable collection of ObjectMetadata, which can be iterated using 'await foreach'. Options can be provided to filter the list.

```csharp
IAsyncEnumerable<ObjectMetadata> ListAsync(NatsObjListOpts? opts = default, CancellationToken cancellationToken = default);
```

--------------------------------

### Java Object Store Operations

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Provides methods for interacting with NATS Object Stores. Supports putting objects from streams, files, or byte arrays, retrieving objects, updating metadata, deleting objects, adding links, sealing buckets, listing objects, and watching for changes.

```Java
/**
     * @return the name
     */
    String getBucketName();

    /**
     * Place the contents of the input stream into a new object.
     */
    ObjectInfo put(ObjectMeta meta, InputStream inputStream) throws IOException, JetStreamApiException, NoSuchAlgorithmException;

    /**
     * Place the contents of the input stream into a new object.
     */
    ObjectInfo put(String objectName, InputStream inputStream) throws IOException, JetStreamApiException, NoSuchAlgorithmException;

    /**
     * Place the bytes into a new object.
     */
    ObjectInfo put(String objectName, byte[] input) throws IOException, JetStreamApiException, NoSuchAlgorithmException;

    /**
     * Place the contents of the file into a new object using the file name as the object name.
     */
    ObjectInfo put(File file) throws IOException, JetStreamApiException, NoSuchAlgorithmException;

    /**
     * Get an object by name from the store, reading it into the output stream, if the object exists.
     */
    ObjectInfo get(String objectName, OutputStream outputStream) throws IOException, JetStreamApiException, InterruptedException, NoSuchAlgorithmException;

    /**
     * Get the info for an object if the object exists / is not deleted.
     */
    ObjectInfo getInfo(String objectName) throws IOException, JetStreamApiException;

    /**
     * Get the info for an object if the object exists, optionally including deleted.
     */
    ObjectInfo getInfo(String objectName, boolean includingDeleted) throws IOException, JetStreamApiException;

    /**
     * Update the metadata of name, description or headers. All other changes are ignored.
     */
    ObjectInfo updateMeta(String objectName, ObjectMeta meta) throws IOException, JetStreamApiException;

    /**
     * Delete the object by name. A No-op if the object is already deleted.
     */
    ObjectInfo delete(String objectName) throws IOException, JetStreamApiException;

    /**
     * Add a link to another object. A link cannot be for another link.
     */
    ObjectInfo addLink(String objectName, ObjectInfo toInfo) throws IOException, JetStreamApiException;

    /**
     * Add a link to another object store (bucket).
     */
    ObjectInfo addBucketLink(String objectName, ObjectStore toStore) throws IOException, JetStreamApiException;

    /**
     * Close (seal) the bucket to changes. The store (bucket) will be read only.
     */
    ObjectStoreStatus seal() throws IOException, JetStreamApiException;

    /**
     * Get a list of all object [infos] in the store.
     */
    List<ObjectInfo> getList() throws IOException, JetStreamApiException, InterruptedException;

    /**
     * Create a watch on the store (bucket).
     */
    NatsObjectStoreWatchSubscription watch(ObjectStoreWatcher watcher, ObjectStoreWatchOption... watchOptions) throws IOException, JetStreamApiException, InterruptedException;

    /**
     * Get the ObjectStoreStatus object.
     */
    ObjectStoreStatus getStatus() throws IOException, JetStreamApiException;
```

--------------------------------

### PUT /object-store/put

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Places content into a new object within the store.

```APIDOC
## PUT /object-store/put

### Description
Uploads content to an object store bucket. Supports streams, byte arrays, or files.

### Method
PUT

### Endpoint
/object-store/put

### Parameters
#### Request Body
- **objectName** (string) - Required - The name of the object.
- **data** (binary) - Required - The content to store.

### Request Example
{
  "objectName": "file.txt",
  "data": "base64_encoded_content"
}

### Response
#### Success Response (200)
- **ObjectInfo** (Object) - Metadata about the stored object.

#### Response Example
{
  "name": "file.txt",
  "size": 1024,
  "checksum": "sha256_hash"
}
```

--------------------------------

### Manage NATS Object Store Buckets

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Methods for retrieving or deleting an object store bucket within the NATS JetStream context.

```csharp
ValueTask<INatsObjStore> GetObjectStoreAsync(string bucket, CancellationToken cancellationToken = default);
ValueTask<bool> DeleteObjectStore(string bucket, CancellationToken cancellationToken);
```

--------------------------------

### Java Object Store Interface Definition

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Basic interface definition for the Object Store management context in Java.

```java
public interface ObjectStore {
    /**
     * Get the name of the object store's bucket.
     */
}
```

--------------------------------

### DELETE /object-store/delete

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Deletes an object from the store by name.

```APIDOC
## DELETE /object-store/delete

### Description
Removes an object from the specified bucket. If the object is already deleted, this is a no-op.

### Method
DELETE

### Endpoint
/object-store/delete

### Parameters
#### Query Parameters
- **objectName** (string) - Required - The name of the object to delete.

### Response
#### Success Response (200)
- **ObjectInfo** (Object) - Metadata of the deleted object.

#### Response Example
{
  "name": "file.txt",
  "deleted": true
}
```

--------------------------------

### Manage Object Metadata and Links

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Methods to update existing object metadata or create links between objects within the store.

```csharp
ValueTask<ObjectMetadata> UpdateMetaAsync(string key, ObjectMetadata meta, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> PutAsync(ObjectMetadata meta, Stream stream, bool leaveOpen = false, CancellationToken cancellationToken = default);
```

--------------------------------

### POST /bucket/seal

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Seals the object store to prevent further modifications.

```APIDOC
## POST /bucket/seal

### Description
Seals the current object store. Once sealed, no further modifications are permitted.

### Method
POST

### Endpoint
/bucket/seal
```

--------------------------------

### Delete Object by Key - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Deletes an object from the store using its unique key. This operation may fail if the object metadata is invalid or if chunks cannot be purged, potentially throwing an NatsObjException.

```csharp
ValueTask DeleteAsync(string key, CancellationToken cancellationToken = default);
```

--------------------------------

### Perform Object Store Operations

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Core interface methods for reading and writing objects to a NATS bucket. Supports both byte array and stream-based data transfer.

```csharp
ValueTask<byte[]> GetBytesAsync(string key, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> GetAsync(string key, Stream stream, bool leaveOpen = false, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> PutAsync(string key, byte[] value, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> PutAsync(string key, Stream stream, bool leaveOpen = false, CancellationToken cancellationToken = default);
```

--------------------------------

### Watch for Store Changes - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Monitors the underlying object store for changes and provides real-time updates on object metadata. It returns an asynchronous enumerable collection of ObjectMetadata, suitable for use with 'await foreach'. Watch options can be specified.

```csharp
IAsyncEnumerable<ObjectMetadata> WatchAsync(NatsObjWatchOpts? opts = default, CancellationToken cancellationToken = default);
```

--------------------------------

### DELETE /objects/{key}

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Deletes an object from the store by its key.

```APIDOC
## DELETE /objects/{key}

### Description
Removes an object from the store using its key.

### Method
DELETE

### Endpoint
/objects/{key}

### Parameters
#### Path Parameters
- **key** (string) - Required - The unique identifier for the object to delete.
```

--------------------------------

### ObjectStoreManager API

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Manages object store instances, including creation, retrieval, update, and deletion.

```APIDOC
## ObjectStoreManager API

### Description
The ObjectStoreManager is used to manage object stores. It provides methods for CRUD operations on object stores.

### Methods

#### ObjectStore

*   **Description**: Looks up and binds to an existing object store instance.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `bucket` (string) - The name of the object store bucket.
*   **Returns**: `ObjectStore` interface and an `error`. Returns `ErrBucketNotFound` if the bucket does not exist.

#### CreateObjectStore

*   **Description**: Creates a new object store with the given configuration.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `cfg` (ObjectStoreConfig) - The configuration for the new object store.
*   **Returns**: `ObjectStore` interface and an `error`. Returns `ErrBucketExists` if the bucket already exists.

#### UpdateObjectStore

*   **Description**: Updates an existing object store with the given configuration.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `cfg` (ObjectStoreConfig) - The updated configuration for the object store.
*   **Returns**: `ObjectStore` interface and an `error`. Returns `ErrBucketNotFound` if the bucket does not exist.

#### CreateOrUpdateObjectStore

*   **Description**: Creates a new object store if it does not exist, or updates an existing one with the given configuration.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `cfg` (ObjectStoreConfig) - The configuration for the object store.
*   **Returns**: `ObjectStore` interface and an `error`.

#### DeleteObjectStore

*   **Description**: Deletes the specified object store.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
    *   `bucket` (string) - The name of the object store bucket to delete.
*   **Returns**: An `error`. Returns `ErrBucketNotFound` if the bucket does not exist.

#### ObjectStoreNames

*   **Description**: Retrieves a list of object store bucket names.
*   **Parameters**: 
    *   `ctx` (context.Context) - The context for the operation.
*   **Returns**: `ObjectStoreNamesLister` which provides a channel to receive bucket names.
```

--------------------------------

### Add Link to Object - C#

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

Adds a link to another object within the NATS Object Store. This method supports linking to another object using its metadata or directly to an object store. It returns the metadata of the newly created link object.

```csharp
ValueTask<ObjectMetadata> AddLinkAsync(string link, string target, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> AddLinkAsync(string link, ObjectMetadata target, CancellationToken cancellationToken = default);
ValueTask<ObjectMetadata> AddBucketLinkAsync(string link, INatsObjStore target, CancellationToken cancellationToken = default);
```

--------------------------------

### Perform Object Operations with Go

Source: https://github.com/nats-io/nats.docs/blob/master/using-nats/developing-with-nats/js/object.md

The ObjectStore interface defines methods for interacting with objects inside a store, such as Put, Get, and Delete. It includes convenience methods for handling byte slices, strings, and files.

```go
type ObjectStore interface {
	Put(ctx context.Context, obj ObjectMeta, reader io.Reader) (*ObjectInfo, error)
	PutBytes(ctx context.Context, name string, data []byte) (*ObjectInfo, error)
	PutString(ctx context.Context, name string, data string) (*ObjectInfo, error)
	PutFile(ctx context.Context, file string) (*ObjectInfo, error)
	Get(ctx context.Context, name string) (ObjectResult, error)
}
```