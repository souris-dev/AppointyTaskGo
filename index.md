# Mock Instagram Backend

This is an API created for a mock Instagram-like application for the internship recruitment task of Appointy in September, 2021 for Vellore Institute of Technology.
This API has been used using Golang and using `net/http` (along with other standard libraries) and no other external libraries or frameworks have been used as per the directions.

## Technologies Used

This backend is written in Golang and uses standard libraries for JSON encoding/decoding and `net/http` for routing and serving. The API is designed to communicate in JSON with its client.
MongoDB is used as the database as specified in the requirements of the task.

## Building the source
Ensure that Golang is installed and set up in your system.
First, clone this repository into a folder. Change directory to that folder.

```sh
git clone <this repository's URL.git>
cd AppointyTaskGo
```

Next, build the project. The API uses the Go driver for MongoDB. Building the project should automatically fetch the dependencies.
From within this directory (project root), run:

```sh
go build
```

## Running the API

To run the API, first set the environment variables `MONGODB_URI`, `MONGODB_DBNAME` and `APPYINSTA_PORT`.

### Linux

```sh
export MONGODB_URI=<your connection string>
export MONGODB_DBNAME=<your database name>
export APPYINST_PORT=<on which port to run the server>
```

### Windows (Powershell)
```ps
$Env:MONGODB_URI = "<...>"
$Env:MONGODB_DBNAME = "<...>"
$Env:APPYINSTA_PORT = "<port>"
```

You can also set these environment variables using other methods.

After this, run the executable created after building it.

### Linux
```sh
./appyinsta
```
You may have to change permissions for executing this.

### Windows (Powershell)
```ps
./appyinsta
```

The server will now run on the specified port (as specified in the `APPYINSTA_PORT` environment variable).

## API Specification

A simple overview of the API is as follows. The API has been designed and created as per the requirements specified in the task.


<table>
  <tr>
    <th>Route</th>
    <th>Method</th>
    <th>Description</th>
    <th>Request Body (Sample)</th>
    <th>Response Body</th>
  </tr>
  
  <tr>
    <td>/users</td>
    <td>POST</td>
    <td>Create a user</td>
    <td>
      <pre>
json
{
  "name": "(name)",
  "email": "(email)",
  "password": "(password)"
}
      </pre>
      The password is hashed again at the server. All fields are compulsory. <br/>
    </td>
    <td>
    <pre>
json
{
  "id": "(user ID)"
}
    </pre>
      The user ID (MongoDB object ID) of the new user is returned after user creation.
    </td>
  </tr>
  <tr>
    <td>/users/&lt;userID&gt;</td>
    <td>GET</td>
    <td>Retrieve information about a user</td>
    <td>N/A</td>
    <td><pre>
json
{
  "id": "(user ID)",
  "name": "(name)",
  "email": "(email)",
}
      </pre>
      The <i>id</i> field has the same user ID as specified in the URL.
    </td>
  </tr>
  <tr>
    <td>/posts</td>
    <td>POST</td>
    <td>Create a post</td>
    <td>
    <pre>
json
{
    "posted_by": "(user ID )",
    "caption": "(caption)",
    "img_url": "(image URL)"
}
    </pre>
      The <i>posted_by</i> field contains the user ID of the user who created this post.
      While post creation, the timestamp of its creation is recorded at the server.
    </td>
    <td>
    <pre>
json
{
  "id": "(post ID)"
}
    </pre>
      The post ID (MongoDB object ID) of the new post is returned after post creation.
    </td>
  </tr>
  <tr>
    <td>/posts/&lt;postID&gt;</td>
    <td>GET</td>
    <td>Retrieve information about a post</td>
    <td>N/A</td>
    <td>
    <pre>
json
{
  "id": "(post ID)",
  "posted_by": "(user ID)",
  "caption": "(caption)",
  "img_url": "(image URL)",
  "posted_on": "(timestamp)"
}
    </pre>
      The <i>id</i> field has the same post ID as specified in the URL.
    </td>
  </tr>
  <tr>
    <td>/posts/users/&lt;userID&gt;</td>
    <td>GET</td>
    <td>Retrieve the posts created by the user, latest first.</td>
    <td>
      <b>First Request</b><br /><br />
      For the first request, the <i>first_request</i> must be set to true. <br />
      The <i>last_id</i> field can be set as the userID (or any post ID). <br />
      The <i>last_posted_on</i> should be any string of a valid timestamp format (for example: 
      2021-10-09T12:17:11.478Z). <br />
      The request body format for the first request is as follows:
      <pre>
json
{
  "last_id": "(user ID)",
  "last_posted_on": "(timestamp)",
  "n_new": 3,
  "first_request": true
}
    </pre>
      The <i>n_new</i> field sets how many posts should be retrieved.
      <br /><br />
      <b>Subsequent Requests</b><br /><br />
       For the subsequent requests, the request format is similar.<br />
       The <i>last_id</i> field should have the post ID of the last post <br />
       of the previous response. <br />
       The <i>last_posted_on</i> field should have the posted_on timestamp of <br />
       the last post of the previous response. <br />
       The <i>first_request</i> field must be set to false, or can be omitted.
    </td>
    <td>
     <pre>
json
[
    {
        "id": "(post ID)",
        "posted_by": "(user ID)",
        "caption": "(caption)",
        "img_url": "(image URL)",
        "posted_on": "(timestamp)"
    },
    ...
    {
        "id": "(post ID)",
        "posted_by": "(user ID)",
        "caption": "(caption)",
        "img_url": "(image URL)",
        "posted_on": "(timestamp)"
    }
]
    </pre>
      This array will contain maximum <i>n_new</i> number of posts (as specified in the request body).
      The posts are returned in the most recent first order.
    </td>
  </tr>
    
</table>


## Running Unit Tests

Ensure that the required environment variables are set (see the "Running the API section"), and go to the project root directory and run:

```sh
go test appyinsta/api/handlers
```

(As of now, the tests depend on existing data in the database, and will fail if that data is not present. A mechanism to add test data will be added.)

## Licence

Will be added soon.



2021, Souris Ash
