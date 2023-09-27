
package edu.ucla.library.validation.handlers;

import static edu.ucla.library.validation.MediaType.APPLICATION_JSON;

import info.freelibrary.util.HTTP;
import java.util.List;

import edu.ucla.library.validation.JsonKeys;

import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpHeaders;
import io.vertx.core.json.JsonObject;
import io.vertx.core.buffer.Buffer;

import io.vertx.ext.web.RoutingContext;
import io.vertx.ext.web.FileUpload;

/**
 * A handler that processes status information requests.
 */
public class UploadHandler implements Handler<RoutingContext> {

    /** The handler's copy of the Vert.x instance. */
    private final Vertx myVertx;

    /**
     * Creates a handler that returns a status response.
     *
     * @param aVertx A Vert.x instance
     */
    public UploadHandler(final Vertx aVertx) {
        myVertx = aVertx;
    }

    @Override
    public void handle(final RoutingContext aContext) {
        List<FileUpload> fileUploads = aContext.fileUploads();
        for (FileUpload fileUpload : fileUploads) {
            String uploadedFileName = fileUpload.fileName();
            Buffer uploadedFile = myVertx.fileSystem().readFileBlocking(fileUpload.uploadedFileName());
            String fileContents = uploadedFile.toString();

            // Send the contents back to the browser
            aContext.response().end(fileContents);
        }

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
