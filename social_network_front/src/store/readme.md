# Redux: 
    Lets us store all the data into one storage unit (usually called store), its like central database for frontend. 
    This store gives us only one place that we need to update, makes data flow transparent and predictable what makes data(state) managment easier.
    Not right tool for every project.

### Pros
- Predictable state changes
- Centralized state
- Easy debugging
- Preserve page state
- Undo/redo
- Ecosystem of add-ons

### Cons
- Complexity (funtional programing princibles)
- Verbosity (needs some boilerplate code)

Reducer(eventHandlers) - 
    Is function that takes a current version of the store and returns an updated version of the store
Reducer -> Action(event) - 
    Action is js object that describes to reducer what should it update.

    Reducers must always follow some specific rules:
    1. They should only calculate the new state value based on the state and action arguments
    2. They are not allowed to modify the existing state. Instead, they must make immutable updates, by copying the existing state and making changes to the copied values.
    3. They must not do any asynchronous logic, calculate random values, or cause other "side effects"

"Slices" in redux means store properties (4 slices)

```
    {
        categories: [],
        products; {],
        cart :{},
        user:{}
    }
```
Dispatch is like entry point to the store. Every action(event) goes through the same point. 

    	``
    	Action(event) --dispatch --> Store ----> Reducer (eventHandler) 
    	                                                ||
    	                                                \/
    	                             Store <---- Reducer
    	``


In redux - functions have to be pure functions
Easiest way to change state is to take obj and spread, after that change values {...OldObject, name : "XXX", age :30, newValue : "ok"}

    ``  
        let obj1 = { id: 1,  myname : "silver" }
        let obj2 = {...obj1, myname: "juku", id: 2, doing: "nothing"}

        obj1 ==> { id: 1,  myname : "silver" }
        obj2 ==> { id: 2,  myname : "juku", doing: "nothing}
    ``
