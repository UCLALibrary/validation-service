
package info.freelibrary.vertx.template.verticles;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import info.freelibrary.util.Constants;
import info.freelibrary.util.HTTP;

import info.freelibrary.vertx.template.Config;
import info.freelibrary.vertx.template.Op;

import io.vertx.config.ConfigRetriever;
import io.vertx.core.DeploymentOptions;
import io.vertx.core.Vertx;
import io.vertx.ext.web.client.WebClient;
import io.vertx.junit5.VertxExtension;
import io.vertx.junit5.VertxTestContext;

/**
 * Tests the main verticle of the Vert.x application.
 */
@ExtendWith(VertxExtension.class)
public class MainVerticleTest {

    /** A Vert.x configuration retriever. */
    private static ConfigRetriever myConfigRetriever;

    /**
     * Sets up the test.
     *
     * @param aVertx A Vert.x instance
     * @param aContext A test context
     */
    @BeforeAll
    public static void setUp(final Vertx aVertx, final VertxTestContext aContext) {
        myConfigRetriever = ConfigRetriever.create(aVertx);

        // Store the relevant test configuration details in our deployment options
        myConfigRetriever.getConfig().onFailure(aContext::failNow).onSuccess(config -> {
            final DeploymentOptions options = new DeploymentOptions().setConfig(config);

            aVertx.deployVerticle(MainVerticle.class.getName(), options).onFailure(aContext::failNow).onSuccess(id -> {
                aContext.completeNow();
            });
        });
    }

    /**
     * Tests the server can start successfully.
     *
     * @param aContext A test context
     */
    @Test
    public void testThatTheServerIsStarted(final Vertx aVertx, final VertxTestContext aContext) {
        myConfigRetriever.getConfig().onFailure(aContext::failNow).onSuccess(config -> {
            final int port = config.getInteger(Config.HTTP_PORT);
            final WebClient client = WebClient.create(aVertx);

            client.get(port, Constants.INADDR_ANY, Op.GET_STATUS.endpoint()).send().onSuccess(response -> {
                aContext.verify(() -> {
                    Assertions.assertEquals(HTTP.OK, response.statusCode());
                }).completeNow();
            }).onFailure(aContext::failNow);
        });
    }
}
