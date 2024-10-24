class AIProvider {
  constructor(name, model, baseURL) {
    this.name = name || "";
    this.model = model || "";
    this.password = "";
    this.baseURL = baseURL || "";
    this.proxyEndpoint = "";
    this.proxyPort = "";
    this.endpointName = "";
    this.engine = "";
    this.temperature = 0.0;
    this.providerRegion = "";
    this.providerId = "";
    this.compartmentId = "";
    this.topP = 0.0;
    this.topK = 0;
    this.maxTokens = 0;
    this.organizationId = "";
    this.customHeaders = [];
  }
}

export default AIProvider;
