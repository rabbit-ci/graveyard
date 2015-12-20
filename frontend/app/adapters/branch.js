import ApplicationAdapter from './application';

export default ApplicationAdapter.extend({
  urlTemplate: '{+host}/projects/{projectId}/branches{/branchId}',

  urlSegments: {
    projectId: function(type, id, snapshot, requestType) {
      return snapshot.get('project').id;
    },

    branchId: function(type, id, snapshot, requestType) {
      return snapshot.id;
    },

    host: function() {
      return this.get('host');
    }
  }
});
