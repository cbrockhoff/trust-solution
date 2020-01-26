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
import org.springframework.web.client.RestTemplate;
import org.springframework.util.MultiValueMap;
import org.springframework.util.LinkedMultiValueMap;
import works.weave.socks.orders.config.RestProxyTemplate;
import org.springframework.http.HttpMethod;

import java.io.IOException;
import java.net.URI;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Future;

//import static org.springframework.hateoas.MediaTypes.HAL_JSON;

@Service
public class AsyncGetService {
    private final Logger LOG = LoggerFactory.getLogger(getClass());

    private final RestProxyTemplate restProxyTemplate;

    private final RestTemplate halTemplate;

    private final static String HAL_JSON = "application/hal+json";

    @Autowired
    public AsyncGetService(RestProxyTemplate restProxyTemplate) {
        this.restProxyTemplate = restProxyTemplate;
        this.halTemplate = new RestTemplate(restProxyTemplate.getRestTemplate().getRequestFactory());

        ObjectMapper objectMapper = new ObjectMapper();
        objectMapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        objectMapper.registerModule(new Jackson2HalModule());
        MappingJackson2HttpMessageConverter halConverter = new MappingJackson2HttpMessageConverter();
        halConverter.setSupportedMediaTypes(Arrays.asList(MediaTypes.HAL_JSON));
        halConverter.setObjectMapper(objectMapper);
        halTemplate.setMessageConverters(Collections.singletonList(halConverter));
    }

    @Async
    public <T> Future<Resource<T>> getResource(URI url, String userToken, String serviceToken, TypeReferences.ResourceType<T> type) throws
            InterruptedException, IOException {

        LinkedMultiValueMap<String, String>headers = new LinkedMultiValueMap<String, String>();
        headers.add("Accept", HAL_JSON);
        headers.add("User-Token", userToken);
        headers.add("Service-Token", serviceToken);

        RequestEntity<Void> request = new RequestEntity<>(headers, HttpMethod.GET, url);

        LOG.debug("Requesting: " + request.toString());
        Resource<T> body = restProxyTemplate.getRestTemplate().exchange(request, type).getBody();
        LOG.debug("Received: " + body.toString());
        return new AsyncResult<>(body);
    }

    @Async
    public <T> Future<Resources<T>> getDataList(URI url, String userToken, String serviceToken, TypeReferences.ResourcesType<T> type) throws
            InterruptedException, IOException {

        LinkedMultiValueMap<String, String> headers = new LinkedMultiValueMap<String, String>();
        headers.add("Accept", HAL_JSON);
        headers.add("User-Token", userToken);
        headers.add("Service-Token", serviceToken);

        RequestEntity<Void> request = new RequestEntity<>(headers, HttpMethod.GET, url);
        LOG.debug("Requesting: " + request.toString());
        Resources<T> body = restProxyTemplate.getRestTemplate().exchange(request, type).getBody();
        LOG.debug("Received: " + body.toString());
        return new AsyncResult<>(body);
    }

    @Async
    public <T> Future<List<T>> getDataList(URI url, String userToken, String serviceToken, ParameterizedTypeReference<List<T>> type) throws
            InterruptedException, IOException {

        LinkedMultiValueMap<String, String> headers = new LinkedMultiValueMap<String, String>();
        headers.add("Accept", "application/json");
        headers.add("User-Token", userToken);
        headers.add("Service-Token", serviceToken);

        RequestEntity<Void> request = new RequestEntity<>(headers, HttpMethod.GET, url);
        LOG.debug("Requesting: " + request.toString());
        List<T> body = restProxyTemplate.getRestTemplate().exchange(request, type).getBody();
        LOG.debug("Received: " + body.toString());
        return new AsyncResult<>(body);
    }

    @Async
    public <T, B> Future<T> postResource(URI uri, B body, String userToken, String serviceToken, ParameterizedTypeReference<T> returnType) {

        LinkedMultiValueMap<String, String> headers = new LinkedMultiValueMap<String, String>();
        headers.add("Accept", "application/json");
        headers.add("Content-Type", "application/json");
        headers.add("User-Token", userToken);
        headers.add("Service-Token", serviceToken);

        //RequestEntity<B> request = RequestEntity.post(uri).contentType(MediaType.APPLICATION_JSON).accept(MediaType.APPLICATION_JSON).body(body);
        RequestEntity<B> request = new RequestEntity<B>(body, headers, HttpMethod.POST, uri);
        LOG.debug("Requesting: " + request.toString());
        T responseBody = restProxyTemplate.getRestTemplate().exchange(request, returnType).getBody();
        LOG.debug("Received: " + responseBody);
        return new AsyncResult<>(responseBody);
    }
}
