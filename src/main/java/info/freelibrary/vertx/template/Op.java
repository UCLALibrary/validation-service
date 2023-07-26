
package info.freelibrary.vertx.template;

/**
 * Template API operation IDs.
 */
public enum Op {

    /** An operation for the status functionality. */
    GET_STATUS("getStatus", "/status");

    /** An operation's ID. */
    private String myID;

    /** An operation's endpoint. */
    private String myEndpoint;

    /**
     * Creates a new operation.
     *
     * @param aID A handler name
     * @param aEndpoint An endpoint
     */
    Op(final String aID, final String aEndpoint) {
        myID = aID;
        myEndpoint = aEndpoint;
    }

    /**
     * Gets the operation's ID.
     *
     * @return An ID
     */
    @SuppressWarnings({ "PMD.ShortMethodName" })
    public String id() {
        return myID;
    }

    /**
     * Gets the operation's endpoint.
     *
     * @return An endpoint
     */
    public String endpoint() {
        return myEndpoint;
    }
}
