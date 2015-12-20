import Ember from 'ember';

export default Ember.Route.extend({
  model: function(params) {
    return this.store.find('branch', params.branchId, {project: params.projectId});
  }
});
