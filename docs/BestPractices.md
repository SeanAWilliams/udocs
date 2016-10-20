# Documentation Best Practices

The following topics cover some best practices to assist you in documenting your application/service. When writing product 
documentation, the most important thing to keep in mind is that your documentation guide is the only lifeline your consumers
 have when trying to use/understand/debug your service. Your customers are unable to `Google` their questions about your
 service, or seek help on `StackOverflow`. If your consumers are unable to find the information they are looking for,
 the only option they are left with is opening a JIRA ticket.


### Which features of my service should be documented?
The features you should document will vary depending on the application/service. However, you may want to consider covering the
following points as they pertain to your product:

* What is my service and what does it do?
* What purpose does my service serve? Which problems is it trying to solve?
* Who are the consumers of my service? Which of their use-case scenarios does my app address?
* What is the SLA ([Service Level Agreement](https://www.paloaltonetworks.com/resources/learning-center/what-is-a-service-level-agreement-sla.html)) between my app and its consumers?
* How can new users get started with consuming my service? (Installation, setup, configuration, examples, etc.)
* How does my service fit into the rest of my product's/platforms'secosystem? Does it depend on, or integrate with, any other apps, services, or technologies?
* How do consumers interact with my service? Is there an API (yes), or any client tools they need to be familiar with?
* Are there any known quirks/bugs my consumers should be aware of?

### What are some common items to include in my documentation?
When deciding on the types of content you want to include in your documentation guide, you should consider including
something about the relevant items:

* How-to guides and tutorials
    * Installation, setup, and configuration settings
    * Examples of how to integrate it with other services
    * Some form of a “Hello, world!” use-case of your service
    * CI pipeline configuration
    * Internal dev/test guides for setting up local environments
* Technical specifications
    * The API of your service
    * Architectural diagrams of your service
    * Code examples
    * Infrastructure requirements for using your service
* Links to external content
    * [GitHub](https://github.com) repo for your service
    * [Swagger](http://swagger.io/getting-started/) endpoint for your service
    * Official documentation guides for the key technologies that your service leverages
    * Related blogs or videos that you think your consumers may find helpful
* Support/Contact info
    * Information about the team/cell who owns this service
    * Change notifications: How will your customers be notified of changes to your service?
    * S.O.P. (Standard Operating Procedure) for common issues with instructions on how to debug them


### What are some items I should NOT include in my documentation?
The following are examples of the types of information you should omit from your documentation guides. This list is not
comprehensive so please use your best judgment when unsure of whether or not to include something in your guide.

* Usernames or passwords
* PEM or cert files
* IP addresses, ports, or any other internal network information


### API documentation
It is important that you clearly document all aspects of your service's API. At minimum, you should enumerate all of
the public methods of your API with instructions on how to call these methods, as well as the keys/fields/parameters of
these methods. For any given key/field/parameter of a method in your API, you should make clear the following:

* Is this field required?
* What is the value type of this field?
* Does this field require a unique value?
* Does this field depend on, or influence any other fields?
* Is this field's value user-specified? If not, what is the set of possible values this field can have?

Once you have clearly documented the methods and parameters of your API, you should include examples of common
use-cases for each of your API's methods. Make sure to include an explanation of what is happening in each of these examples,
and the parameters used in those methods. You should also provide your customers a way to experiment with using
your API. There is no required tool for this, but [Swagger](http://swagger.io/getting-started/) is a good option. Make sure
to sandbox any environment you are letting your consumers use to test your API. If any of this seems tedious or unnecessary, remember
all of the times you have tried to debug some service or program you are using, only to find its documentation completely lacking.

