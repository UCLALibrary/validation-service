
package info.freelibrary.vertx.template.handlers;

import static info.freelibrary.vertx.template.MediaType.APPLICATION_JSON;

import info.freelibrary.util.HTTP;

import info.freelibrary.vertx.template.JsonKeys;

import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpHeaders;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;

/**
 * A handler that processes status information requests.
 */
public class StatusHandler implements Handler<RoutingContext> {

    /** The handler's copy of the Vert.x instance. */
    private final Vertx myVertx;

    /**
     * Creates a handler that returns a status response.
     *
     * @param aVertx A Vert.x instance
     */
    public StatusHandler(final Vertx aVertx) {
        myVertx = aVertx;
    }

    @Override
    public void handle(final RoutingContext aContext) {
        aContext.response().setStatusCode(HTTP.OK).putHeader(HttpHeaders.CONTENT_TYPE, APPLICATION_JSON.toString())
                .end(new JsonObject().put(JsonKeys.STATUS, "ok").encodePrettily());
    }

    /**
     * Gets the Vert.x instance associated with this handler.
     *
     * @return The Vert.x instance associated with this handler
     */
    public Vertx getVertx() {
        return myVertx;
    }
}
