package works.weave.socks.orders.services;

import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.hateoas.MediaTypes;
import org.springframework.hateoas.Resource;
import org.springframework.hateoas.Resources;
import org.springframework.hateoas.hal.Jackson2HalModule;
import org.springframework.hateoas.mvc.TypeReferences;
import org.springframework.http.MediaType;
import org.springframework.http.RequestEntity;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.scheduling.annotation.Async;
import org.springframework.scheduling.annotation.AsyncResult;
import org.springframework.stereotype.Service;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.web.client.RestTemplate;
import org.springframework.hateoas.mvc.TypeReferences;
import works.weave.socks.orders.config.RestProxyTemplate;

import java.util.concurrent.CompletionException;
import java.util.Base64;
import java.net.http.HttpClient;
import java.net.http.HttpHeaders;
import java.net.http.HttpRequest;

import java.net.http.HttpResponse.BodyHandlers;
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpResponse;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.concurrent.Future;
import java.util.concurrent.CompletableFuture;


@Service
public class TokenGetService {
    private final Logger LOG = LoggerFactory.getLogger(getClass());

    private final static Object tokenLock = new Object();
    private static String token = "";
    private static long getNewTokenAt = 0;

    public static class TokenResponse {
        public String Token;

        public TokenResponse() {
            this.Token = "None";
        }

        public TokenResponse(String token) {
            this.Token = token;
        }

        public void setToken(String token) {
            this.Token = token;
        }
    }

    public String getServiceToken() {

        synchronized (tokenLock) {
            if (getNewTokenAt - System.currentTimeMillis() <= 0) {
                LOG.info("Fetching service token...");
                token = requestToken();
                getNewTokenAt = System.currentTimeMillis() + 4 * 60 * 1000; // Should really get this from the token itself. Oh well.
                LOG.info("Service token: " + token);
            }
        }

        return token;
    }

    private String requestToken() {
        UncheckedObjectMapper objectMapper = new UncheckedObjectMapper();
        String base64EncodedAuthorization;
        try {
            base64EncodedAuthorization = Base64.getEncoder().encodeToString("orders:orders_password".getBytes("utf-8"));
        }
        catch(java.io.UnsupportedEncodingException e){
            base64EncodedAuthorization = "";
        }

        HttpRequest request = HttpRequest.newBuilder(URI.create("http://tokens/tokens"))
              .headers("Accept", "application/json", "Authorization", "Basic " + base64EncodedAuthorization)
              .build();

        try {
            return HttpClient.newHttpClient()
                  .sendAsync(request, BodyHandlers.ofString())
                  .thenApply(HttpResponse::body)
                  .thenApply(objectMapper::readValue).get().Token;
        } catch(java.lang.InterruptedException e) {
            System.out.println(e);
            return "FUCK";
        }
        catch(java.util.concurrent.ExecutionException e) {
            System.out.println(e);
            return "FUCK";
        }
    }

    class UncheckedObjectMapper extends ObjectMapper {
            public TokenResponse readValue(String content) {
                try {
                    return readValue(content, TokenGetService.TokenResponse.class);
                } catch (IOException ioe) {
                    throw new CompletionException(ioe);
                }
            }

    }
}
