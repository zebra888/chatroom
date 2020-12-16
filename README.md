**Purpose of this exercise**   

I'd like evaluate GO buildin /net/rpc package before trying out grpc

**Usage**   

-- start server   
go run server/server.go   

-- start client and login chatter   
go run client/client.go chatter_name   

**Basic functionality**

This is a chat room project that supports 4 basic actions: Login, Logout, Post, and Listen.   
Login and Logout only need chatter name provided.    
Post will post a message to the server along with chatter name. The message will be broadcasted to all other logged in chatters.    
Listen runs on the background once a chatter logs in. It retrive messages posted by other chatters.   

**Conclusion about RPC package**   

Generally easy to use.   
Does not support streaming (at least afaik).  

**TUI**   

Used "github.com/marcusolsson/tui-go" package
