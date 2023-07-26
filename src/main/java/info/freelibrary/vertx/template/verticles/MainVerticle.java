
package info.freelibrary.vertx.template.verticles;

import static info.freelibrary.util.Constants.INADDR_ANY;

import java.io.File;

import info.freelibrary.util.Logger;
import info.freelibrary.util.LoggerFactory;

import info.freelibrary.vertx.template.Config;
import info.freelibrary.vertx.template.MessageCodes;
import info.freelibrary.vertx.template.Op;
import info.freelibrary.vertx.template.handlers.StatusHandler;

import io.vertx.config.ConfigRetriever;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Promise;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerOptions;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.openapi.RouterBuilder;

/**
 * Main verticle that starts the application.
 */
public class MainVerticle extends AbstractVerticle {

    /**
     * A logger for the main verticle.
     */
    private static final Logger LOGGER = LoggerFactory.getLogger(MainVerticle.class, MessageCodes.BUNDLE);

    /**
     * An OpenAPI definition that the main verticle users to route requests.
     */
    private static final String API_SPEC = "src/main/resources/template.yaml";

    /**
     * The main verticle's HTTP server.
     */
    private HttpServer myServer;

    @Override
    public void start(final Promise<Void> aPromise) {
        ConfigRetriever.create(vertx).getConfig().onFailure(aPromise::fail)
                .onSuccess(config -> configureServer(config.mergeIn(config()), aPromise));
    }

    @Override
    public void stop(final Promise<Void> aPromise) {
        myServer.close().onFailure(aPromise::fail).onSuccess(result -> aPromise.complete());
    }

    /**
     * Configure the application server.
     *
     * @param aConfig A JSON configuration
     * @param aPromise A startup promise
     */
    private void configureServer(final JsonObject aConfig, final Promise<Void> aPromise) {
        final String host = aConfig.getString(Config.HTTP_HOST, INADDR_ANY);
        final int port = aConfig.getInteger(Config.HTTP_PORT, 8888);

        RouterBuilder.create(vertx, getRouterSpec()).onFailure(aPromise::fail).onSuccess(routeBuilder -> {
            final HttpServerOptions serverOptions = new HttpServerOptions().setPort(port).setHost(host);

            // Associate handlers with operation IDs from the application's OpenAPI specification
            routeBuilder.operation(Op.GET_STATUS.id()).handler(new StatusHandler(getVertx()));

            myServer = getVertx().createHttpServer(serverOptions).requestHandler(routeBuilder.createRouter());
            myServer.listen().onFailure(aPromise::fail).onSuccess(result -> {
                LOGGER.info(MessageCodes.CODE_001, port);
                aPromise.complete();
            });
        });
    }

    /**
     * Gets the OpenAPI specification used to configure the application's router. If the file doesn't exist on the file
     * system, we assume we're running from a Jar and that it can be found in the classpath.
     *
     * @return The OpenAPI router specification
     */
    private String getRouterSpec() {
        final File specFile = new File(API_SPEC);
        return specFile.exists() ? API_SPEC : specFile.getName();
    }

    /**
     * Starts up the main verticle.
     *
     * @param aArgsArray An array of arguments
     */
    @SuppressWarnings("UncommentedMain")
    public static void main(final String[] aArgsArray) {
        Vertx.vertx().deployVerticle(new MainVerticle());
    }
}
