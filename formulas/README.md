### rules
* one value for each service must have _total postfix (it will be used as the value for request count and request response time)
* the value_name for each item must have the service name followed by an underline as a prefix 
    * there must bet both ```servicename_valuename``` (this is the pattern), can't ```have service``` only
* only valid operators in formula are: ```+,-,(,)```
