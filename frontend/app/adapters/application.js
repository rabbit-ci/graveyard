import DS from 'ember-data';
import UrlTemplates from "ember-data-url-templates";

export default DS.RESTAdapter.extend(UrlTemplates, {
  host: "http://localhost:4000",
  namespace: ""
});
