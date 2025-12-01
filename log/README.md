This package provides a common json logger for our go microservices.

# Logging
There are two ways to use the logger - bare function calls, which use the inbuilt defaultLogger, or building your own logger with New() and passing that around.  New lets you specify the writer, whereas default always uses stdout.

For most of our use cases, using the defaultLogger should be fine.  It's all mutexed and thread safe.

# Attributes
If you want to add custom attributes to your logger, use the With() function.  The format for argument is key, value, key, value and so on.  For example - 

    log.With("serverID", 1234).Error("this is an error)

The above example just has one key value pair, but you can add as many as you'd like at once.

Like the Tracing and Customer example below, With() returns a new logger.  So if you want those attributes to be persistent across logs, capture the logger - 

    myLogger := log.With("serverID", 1234)
    myLogger.Error("error 1")
    myLogger.Error("error 2")
    myLogger.Error("error 3")

In the above, all three errors will have the server ID.

# Tracing

This is shorthand for .With("traceId", "my-trace-id-here").  It's mainly so the user doesn't have to remember, use, or misspell the key.

When we have a trace id we want logged, use the WithTrace() function.  This actually returns a new logger with the trace id inbuilt, and can be used inline OR replace your existing logger.

The trace id gets stored in json field 'traceId'

For example, calling 

    log.WithTrace("12345").Error("this is an error") 

works fine for a one shot.  

However, if you want to log everything with a trace, it may be preferable to just capture the object and reuse it, such as

    myLogger := log.WithTrace("12345")
    myLogger.Error("error 1")
    myLogger.Error("error 2")
    myLogger.Error("error 3")

In the above, all three errors will have the trace ID.

# Customer

This is shorthand for .With("customerId", "my-customer-id-here").  It's mainly so the user doesn't have to remember, use, or misspell the key.

When we have a trace id we want logged, use the WithCustomer() function.  This actually returns a new logger with the trace id inbuilt, and can be used inline OR replace your existing logger.

The trace id gets stored in json field 'customerId'

For example, calling 

    log.WithCustomer("customerX").Error("this is an error") 

works fine for a one shot.  

However, if you want to log everything with a customer, it may be preferable to just capture the object and reuse it, such as

    myLogger := log.WithCustomer("customerX")
    myLogger.Error("error 1")
    myLogger.Error("error 2")
    myLogger.Error("error 3")

In the above, all three errors will have the customer ID.
 
 # Levels
 The four levels are DEBUG, INFO, WARN, and ERROR, in that order of severity

 By default, a new logger sets the logging level to INFO, which means everything INFO and above gets logged, and DEBUGs are dropped.

 If you want to change the log level globally, use 

     log.SetMinLogLevel(log.MinLevelDebug) 

for example to set it to DEBUG globally.  

This can also be called on your logging object, like

     myLogger.SetMinLogLevel(log.MinLevelDebug)